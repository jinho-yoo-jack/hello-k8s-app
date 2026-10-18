// 커밋이 올라오면: 테스트 → 이미지 빌드 → GHCR 푸시 → GitOps 저장소의 이미지 태그 변경
// 실제 배포(클러스터 반영)는 Jenkins 가 하지 않는다. Argo CD 가 GitOps 저장소를 보고 한다.
pipeline {
    agent any

    environment {
        GITHUB_USER = 'jinho-yoo-jack'                                   // 자신의 GitHub 계정(소문자)으로 바꾼다
        IMAGE       = "ghcr.io/${GITHUB_USER}/hello-k8s-app"
        GITOPS_REPO = "github.com/${GITHUB_USER}/hello-k8s-gitops.git"
    }

    triggers {
        pollSCM('H/2 * * * *')   // 2분마다 저장소를 확인해 새 커밋이 있으면 빌드 (로컬 Jenkins 는 GitHub 웹훅을 받을 수 없다)
    }

    options {
        disableConcurrentBuilds()        // 같은 잡이 동시에 두 번 돌며 GitOps 저장소에 충돌 커밋을 내지 않도록
        buildDiscarder(logRotator(numToKeepStr: '20'))
    }

    stages {
        stage('Test') {
            steps {
                sh 'docker build --target test -t hello-k8s-app:test .'
            }
        }

        stage('Build') {
            steps {
                script {
                    // 태그 예: 7-3f9c2ab  (빌드 번호 - 커밋 해시 앞 7자리) → 순서와 출처를 둘 다 알 수 있다
                    env.TAG = "${env.BUILD_NUMBER}-${env.GIT_COMMIT.take(7)}"
                }
                sh '''
                    docker build \
                      --build-arg VERSION=$TAG \
                      --label org.opencontainers.image.source=https://github.com/$GITHUB_USER/hello-k8s-app \
                      -t $IMAGE:$TAG .
                '''
            }
        }

        stage('Push') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'github', usernameVariable: 'GH_USER', passwordVariable: 'GH_TOKEN')]) {
                    sh '''
                        echo "$GH_TOKEN" | docker login ghcr.io -u "$GH_USER" --password-stdin
                        docker push $IMAGE:$TAG
                    '''
                }
            }
        }

        stage('Update GitOps') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'github', usernameVariable: 'GH_USER', passwordVariable: 'GH_TOKEN')]) {
                    sh '''
                        rm -rf gitops
                        git clone https://$GH_USER:$GH_TOKEN@$GITOPS_REPO gitops
                        cd gitops
                        yq -i ".image.repository = \\"$IMAGE\\" | .image.tag = \\"$TAG\\"" hello/values.yaml
                        git config user.name  "jenkins"
                        git config user.email "jenkins@k8s-lab.local"
                        git commit -am "deploy: hello-k8s-app $TAG"
                        git push origin HEAD:main
                    '''
                }
            }
        }
    }

    post {
        always {
            sh 'docker logout ghcr.io || true'
        }
        success {
            echo "이미지 ${IMAGE}:${TAG} 를 GitOps 저장소에 반영했습니다. Argo CD 가 곧 배포합니다."
        }
    }
}

pipeline {
    agent any

    environment {
        TARGET_BRANCH = "main"
        GIT_CREDENTIALS_ID = "github-creds"
        ADMIN_EMAIL = "rahul.singhh.144@gmail.com"
    }

    options {
        disableConcurrentBuilds()
        timestamps()
    }

    stages {
        stage('Info') {
            steps {
                script {
                    echo "Building branch: ${env.BRANCH_NAME}"
                }
            }
        }

        stage('Checkout') {
            when {
                expression { env.BRANCH_NAME.startsWith("feat/") }
            }
            steps {
                checkout scm
            }
        }

        stage('Run tests') {
            when {
                expression { env.BRANCH_NAME.startsWith("feat/") }
            }
            steps {
                sh 'go test ./...'
            }
        }

        stage('Attempt merge into target branch') {
            when {
                expression { env.BRANCH_NAME.startsWith("feat/") }
            }
            steps {
                script {
                    sh """
                      git config user.name "jenkins-bot"
                      git config user.email "jenkins-bot@example.com"

                      git fetch origin

                      git checkout -B ${TARGET_BRANCH} origin/${TARGET_BRANCH}

                      set +e
                      git merge --no-commit --no-ff origin/${env.BRANCH_NAME}
                      MERGE_STATUS=\$?
                      set -e

                      if [ "\$MERGE_STATUS" -ne 0 ]; then
                        echo "Conflict detected!"
                        git merge --abort || true
                        exit 99
                      fi

                      git commit -m "auto merge ${env.BRANCH_NAME} into ${TARGET_BRANCH}"
                    """
                }
            }
        }

        stage('Push merged target branch') {
            when {
                expression { env.BRANCH_NAME.startsWith("feat/") }
            }
            steps {
                withCredentials([usernamePassword(credentialsId: env.GIT_CREDENTIALS_ID,
                                                 usernameVariable: 'GIT_USER',
                                                 passwordVariable: 'GIT_PASS')]) {
                    sh """
                      git push https://$GIT_USER:$GIT_PASS@github.com/YOUR_USER/go_cicd_demo.git ${TARGET_BRANCH}
                    """
                }
            }
        }
    }

    post {
        failure {
            script {
                emailext(
                    to: "${ADMIN_EMAIL}",
                    subject: "Merge failed for ${env.BRANCH_NAME}",
                    body: "Merge conflict or error occurred."
                )
            }
        }
    }
}

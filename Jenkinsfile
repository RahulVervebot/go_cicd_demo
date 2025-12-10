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
                    // BRANCH_NAME for multibranch, GIT_BRANCH for classic
                    def branch = env.BRANCH_NAME ?: env.GIT_BRANCH ?: "unknown"
                    // GIT_BRANCH is often like 'origin/main', so we can clean it:
                    branch = branch.replaceFirst(/^origin\//, "")
                    echo "Building branch: ${branch}"
                    // Optionally store cleaned branch in env for later
                    env.EFFECTIVE_BRANCH = branch
                }
            }
        }

        stage('Checkout') {
            when {
                expression { 
                    def branch = env.BRANCH_NAME ?: env.GIT_BRANCH ?: ""
                    branch = branch.replaceFirst(/^origin\//, "")
                    // Only run for feature branches:
                    return branch.startsWith("feat/")
                }
            }
            steps {
                checkout scm
            }
        }

        stage('Run tests') {
            when {
                expression { 
                    def branch = env.BRANCH_NAME ?: env.GIT_BRANCH ?: ""
                    branch = branch.replaceFirst(/^origin\//, "")
                    branch.startsWith("feat/")
                }
            }
            steps {
                sh 'go test ./...'
            }
        }

        stage('Attempt merge into target branch') {
            when {
                expression { 
                    def branch = env.BRANCH_NAME ?: env.GIT_BRANCH ?: ""
                    branch = branch.replaceFirst(/^origin\//, "")
                    branch.startsWith("feat/")
                }
            }
            steps {
                script {
                    sh """
                      git config user.name "jenkins-bot"
                      git config user.email "jenkins-bot@example.com"

                      git fetch origin

                      git checkout -B ${TARGET_BRANCH} origin/${TARGET_BRANCH}

                      set +e
                      git merge --no-commit --no-ff origin/${env.EFFECTIVE_BRANCH}
                      MERGE_STATUS=\$?
                      set -e

                      if [ "\$MERGE_STATUS" -ne 0 ]; then
                        echo "Conflict detected!"
                        git merge --abort || true
                        exit 99
                      fi

                      git commit -m "auto merge ${env.EFFECTIVE_BRANCH} into ${TARGET_BRANCH}"
                    """
                }
            }
        }

        stage('Push merged target branch') {
            when {
                expression { 
                    def branch = env.BRANCH_NAME ?: env.GIT_BRANCH ?: ""
                    branch = branch.replaceFirst(/^origin\//, "")
                    branch.startsWith("feat/")
                }
            }
            steps {
                withCredentials([usernamePassword(credentialsId: env.GIT_CREDENTIALS_ID,
                                                 usernameVariable: 'GIT_USER',
                                                 passwordVariable: 'GIT_PASS')]) {
                    sh """
                      git push https://$GIT_USER:$GIT_PASS@github.com/RahulVervebot/go_cicd_demo.git ${TARGET_BRANCH}
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
                    subject: "Merge failed for branch ${env.EFFECTIVE_BRANCH ?: (env.BRANCH_NAME ?: env.GIT_BRANCH)}",
                    body: "Merge conflict or error occurred. Check job ${env.JOB_NAME} build #${env.BUILD_NUMBER}"
                )
            }
        }
    }
    
}

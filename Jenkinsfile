pipeline {
    agent any

    environment {
        // change to your target integration branch
        TARGET_BRANCH = "develop"
        GIT_CREDENTIALS_ID = "github-creds"  // configure in Jenkins
        ADMIN_EMAIL = "rahul.singhh.144@gmail.com"
    }

    options {
        disableConcurrentBuilds()
        timestamps()
    }

    triggers {
        // For GitHub webhooks, you can also use: pollSCM('* * * * *') OFF
        // Usually, Multibranch + webhook is enough, no need to define triggers here.
    }

    stages {
        stage('Info') {
            steps {
                script {
                    echo "Building branch: ${env.BRANCH_NAME}"
                    if (!env.BRANCH_NAME.startsWith("feat/")) {
                        echo "Not a feature branch, skipping auto-merge."
                    }
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

                      # make sure we have all branches
                      git fetch origin

                      # create a local tracking branch for target
                      git checkout -B ${TARGET_BRANCH} origin/${TARGET_BRANCH}

                      # try merge feature branch into target
                      set +e
                      git merge --no-ff --no-commit origin/${env.BRANCH_NAME}
                      MERGE_STATUS=\$?
                      set -e

                      if [ "\$MERGE_STATUS" -ne 0 ]; then
                        echo "Merge conflict detected between ${env.BRANCH_NAME} and ${TARGET_BRANCH}"
                        # abort merge to keep workspace clean
                        git merge --abort || true
                        exit 99
                      fi

                      echo "Merge clean, committing..."
                      git commit -m "chore: auto-merge ${env.BRANCH_NAME} into ${TARGET_BRANCH} [ci]"
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
                      # push merged ${TARGET_BRANCH} back to origin
                      git push https://$GIT_USER:$GIT_PASS@github.com/YOUR_USER/go_cicd_demo.git ${TARGET_BRANCH}
                    """
                }
            }
        }

        // Optional: Deploy per-branch to different domains/environments
        stage('Deploy to environment') {
            when {
                expression { env.BRANCH_NAME.startsWith("feat/") }
            }
            steps {
                script {
                    echo "Here you could deploy ${TARGET_BRANCH} to some environment."
                    // example:
                    // sh './deploy_to_staging.sh'
                }
            }
        }
    }

    post {
        success {
            script {
                if (env.BRANCH_NAME.startsWith("feat/")) {
                    echo "Auto-merge of ${env.BRANCH_NAME} into ${TARGET_BRANCH} succeeded."
                    // optional: send success email/slack
                }
            }
        }

        failure {
            script {
                if (currentBuild.result == 'FAILURE' || currentBuild.result == 'UNSTABLE') {
                    echo "Build or tests failed for ${env.BRANCH_NAME}."
                }
            }
        }

        // This handles merge conflict specifically (exit 99 in merge stage)
        unstable {
            echo "Build marked unstable."
        }

        always {
            script {
                if (currentBuild.currentResult == 'FAILURE' || currentBuild.currentResult == 'UNSTABLE') {
                    echo "Sending notification about failure/conflict..."

                    // Example: email notification
                    emailext(
                        to: "dev-team@example.com, ${ADMIN_EMAIL}",
                        subject: "Jenkins: Issue while merging ${env.BRANCH_NAME} into ${TARGET_BRANCH}",
                        body: """\
Hello,

There was an issue while processing branch: ${env.BRANCH_NAME}
Job: ${env.JOB_NAME}
Build: ${env.BUILD_NUMBER}
Result: ${currentBuild.currentResult}

Please check Jenkins console logs for details.

Thanks,
Jenkins
"""
                    )
                }
            }
        }
    }
}

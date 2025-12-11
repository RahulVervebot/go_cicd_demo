pipeline {

    agent any

    environment {
        TARGET_BRANCH     = "main"
        GIT_CREDENTIALS_ID = "github-creds"
        ADMIN_EMAIL       = "rahul.singhh.144@gmail.com"
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
                    // Clean up origin/main style
                    branch = branch.replaceFirst(/^origin\//, "")
                    echo "Building branch: ${branch}"
                    env.EFFECTIVE_BRANCH = branch
                }
            }
        }

        stage('Checkout') {
            when {
                expression {
                    def branch = env.BRANCH_NAME ?: env.GIT_BRANCH ?: ""
                    branch = branch.replaceFirst(/^origin\//, "")
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
                    return branch.startsWith("feat/")
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
                    return branch.startsWith("feat/")
                }
            }
            steps {
                script {
                    sh """
                      git config user.name "jenkins-bot"
                      git config user.email "jenkins-bot@example.com"

                      git fetch origin

                      # Checkout latest target branch
                      git checkout -B ${TARGET_BRANCH} origin/${TARGET_BRANCH}

                      # Try merge and capture output
                      set +e
                      git merge --no-commit --no-ff origin/${env.EFFECTIVE_BRANCH} > merge_output.txt 2>&1
                      MERGE_STATUS=\$?
                      set -e

                      if [ "\$MERGE_STATUS" -ne 0 ]; then
                        echo "Conflict detected!"
                        git merge --abort || true
                        exit 99   # non-zero so Jenkins marks build as FAILURE
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
                    return branch.startsWith("feat/")
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

        // Only send success email if it was a feat/* branch
        success {
            script {
                def branch = env.EFFECTIVE_BRANCH ?: (env.BRANCH_NAME ?: env.GIT_BRANCH ?: "unknown")
                branch = branch.replaceFirst(/^origin\//, "")
                if (!branch.startsWith("feat/")) {
                    return
                }

                emailext(
                    to: "${ADMIN_EMAIL}",
                    subject: "Auto-merge SUCCESS: ${branch} -> ${TARGET_BRANCH}",
                    body: """\
Auto-merge completed successfully.

Feature branch: ${branch}
Target branch:  ${TARGET_BRANCH}

Job:   ${env.JOB_NAME}
Build: #${env.BUILD_NUMBER}

The changes from ${branch} have been merged into ${TARGET_BRANCH} and pushed to GitHub.
"""
                )
            }
        }

        failure {
            script {
                def branch = env.EFFECTIVE_BRANCH ?: (env.BRANCH_NAME ?: env.GIT_BRANCH ?: "unknown")
                branch = branch.replaceFirst(/^origin\//, "")
                if (!branch.startsWith("feat/")) {
                    return
                }

                // Capture merge error / conflict details, if available
                def mergeOutput = ""
                if (fileExists('merge_output.txt')) {
                    mergeOutput = readFile('merge_output.txt')
                }

                // Get last commit author's email (the user who pushed)
                def authorEmail = sh(
                    script: "git log -1 --pretty=format:'%ae' || echo ''",
                    returnStdout: true
                ).trim()

                // Send to admin + committer if we have their email
                def recipients = "${ADMIN_EMAIL}"
                if (authorEmail) {
                    recipients = "${recipients}, ${authorEmail}"
                }

                emailext(
                    to: recipients,
                    subject: "Auto-merge FAILED for branch ${branch}",
                    body: """\
Auto-merge failed for:

Feature branch: ${branch}
Target branch:  ${TARGET_BRANCH}

Job:   ${env.JOB_NAME}
Build: #${env.BUILD_NUMBER}

Reason (git output):

${mergeOutput}
"""
                )
            }
        }
    }

}

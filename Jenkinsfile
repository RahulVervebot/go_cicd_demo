pipeline {
  agent any
  
  environment {
    TARGET_BRANCH      = "main"
    GIT_CREDENTIALS_ID = "6f9dfd84-7cc3-412e-8c73-b6c8bd1a3291"
    ADMIN_EMAIL        = "rahul.singhh.144@gmail.com"
    REPO_URL           = "https://github.com/RahulVervebot/go_cicd_demo.git"
    FROM_EMAIL         = "" // optional: set same as your SMTP login, e.g. "yourgmail@gmail.com"
  }

  options {
    disableConcurrentBuilds()
    timestamps()
  }

  stages {

    stage('Info') {
      steps {
        script {
          def branch = env.BRANCH_NAME ?: env.GIT_BRANCH ?: "unknown"
          branch = branch.replaceFirst(/^origin\//, "")
          env.EFFECTIVE_BRANCH = branch
          echo "Building branch: ${env.EFFECTIVE_BRANCH}"
        }
      }
    }

    stage('Checkout') {
      when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
      steps {
        checkout scm
      }
    }

  stage('Capture developer email') {
  when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
  steps {
    script {
      // make sure remote refs exist
      sh "git fetch origin +refs/heads/*:refs/remotes/origin/*"

      // email from the latest commit on the feature branch (NOT from main / NOT from jenkins merge commit)
      def devEmail = sh(
        script: "git show -s --format=%ae origin/${env.EFFECTIVE_BRANCH} || true",
        returnStdout: true
      ).trim()

      // fallback to committer email if author email is blank
      if (!devEmail) {
        devEmail = sh(
          script: "git show -s --format=%ce origin/${env.EFFECTIVE_BRANCH} || true",
          returnStdout: true
        ).trim()
      }

      env.DEVELOPER_EMAIL = devEmail
      echo "Developer email captured from origin/${env.EFFECTIVE_BRANCH}: ${env.DEVELOPER_EMAIL}"
    }
  }
}


    stage('Run tests') {
      when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
      steps {
        sh '''#!/bin/bash -eo pipefail
          if [ -f go.mod ]; then
            echo "go.mod found; running: go test ./..."
            rm -f test_output.txt
            go test ./... | tee test_output.txt
          else
            echo "go.mod not found; skipping tests (recommended: add go.mod and enable go test ./...)"
          fi
        '''
      }
    }

    // If you want MERGE-only (recommended), tell me and I’ll replace these stages.
    stage('Rebase feature onto latest main') {
      when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
      steps {
        script {
          int status = sh(
            script: """#!/bin/bash -e
git config user.name "jenkins-bot"
git config user.email "jenkins-bot@example.com"

git fetch origin

git checkout -B ${env.EFFECTIVE_BRANCH} origin/${env.EFFECTIVE_BRANCH}

set +e
git rebase origin/${TARGET_BRANCH} > rebase_output.txt 2>&1
REBASE_STATUS=\$?
set -e

if [ "\$REBASE_STATUS" -ne 0 ]; then
  echo "Rebase conflict detected!"
  git rebase --abort || true
  exit 98
fi
""",
            returnStatus: true
          )

          if (status == 98) {
            error("Rebase conflict detected between ${env.EFFECTIVE_BRANCH} and ${TARGET_BRANCH}. See rebase_output.txt for details.")
          } else if (status != 0) {
            error("Rebase stage failed with exit code ${status}")
          }
        }
      }
    }

    stage('Merge into main') {
      when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
      steps {
        script {
          int status = sh(
            script: """#!/bin/bash -e
git fetch origin

git checkout -B ${TARGET_BRANCH} origin/${TARGET_BRANCH}

set +e
git merge --no-ff ${env.EFFECTIVE_BRANCH} > merge_output.txt 2>&1
MERGE_STATUS=\$?
set -e

if [ "\$MERGE_STATUS" -ne 0 ]; then
  echo "Merge conflict detected!"
  git merge --abort || true
  exit 99
fi

echo "Merge commit:"
git log -1 --oneline
""",
            returnStatus: true
          )

          if (status == 99) {
            error("Merge conflict detected when merging ${env.EFFECTIVE_BRANCH} into ${TARGET_BRANCH}. See merge_output.txt for details.")
          } else if (status != 0) {
            error("Merge stage failed with exit code ${status}")
          }
        }
      }
    }

    stage('Push main') {
      when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
      steps {
        withCredentials([usernamePassword(
          credentialsId: env.GIT_CREDENTIALS_ID,
          usernameVariable: 'GIT_USER',
          passwordVariable: 'GIT_PASS'
        )]) {
          sh '''
            set -e

            cat > .git_askpass.sh <<'EOF'
#!/bin/sh
case "$1" in
  *Username*) echo "$GIT_USER" ;;
  *Password*) echo "$GIT_PASS" ;;
  *) echo "" ;;
esac
EOF
            chmod +x .git_askpass.sh

            export GIT_ASKPASS="$PWD/.git_askpass.sh"
            export GIT_TERMINAL_PROMPT=0

            git remote set-url origin "$REPO_URL"
            git push origin main

            rm -f .git_askpass.sh
          '''
        }
      }
    }
  }

  post {

    success {
      script {
        def branch = (env.EFFECTIVE_BRANCH ?: env.BRANCH_NAME ?: env.GIT_BRANCH ?: "unknown")
        branch = branch.replaceFirst(/^origin\//, "")

        // Developer email (last commit author)
        def authorEmail = ""
        try {
          authorEmail = sh(script: "git log -1 --pretty=format:%ae || true", returnStdout: true).trim()
        } catch (e) {
          authorEmail = ""
        }

        def recipients = "${ADMIN_EMAIL}"
        if (authorEmail) recipients = "${recipients}, ${authorEmail}"

        def bodyText = """Build & merge successful.

Feature branch: ${branch}
Target branch:  ${TARGET_BRANCH}

Job:   ${env.JOB_NAME}
Build: #${env.BUILD_NUMBER}

Console: ${env.BUILD_URL}console
"""

        def mailArgs = [
          to: recipients,
          subject: "Jenkins SUCCESS: ${branch} merged into ${TARGET_BRANCH}",
          body: bodyText
        ]

        if (env.FROM_EMAIL?.trim()) {
          mailArgs.from = env.FROM_EMAIL.trim()
        }

        emailext(mailArgs)
      }
    }

    failure {
      script {
        def branch = (env.EFFECTIVE_BRANCH ?: env.BRANCH_NAME ?: env.GIT_BRANCH ?: "unknown")
        branch = branch.replaceFirst(/^origin\//, "")

        def mergeOutput  = fileExists('merge_output.txt')  ? readFile('merge_output.txt')  : ""
        def rebaseOutput = fileExists('rebase_output.txt') ? readFile('rebase_output.txt') : ""
        def testOutput   = fileExists('test_output.txt')   ? readFile('test_output.txt')   : ""

        def authorEmail = ""
        try {
          authorEmail = sh(script: "git log -1 --pretty=format:%ae || true", returnStdout: true).trim()
        } catch (e) {
          authorEmail = ""
        }

        def recipients = "${ADMIN_EMAIL}"
        if (authorEmail) recipients = "${recipients}, ${authorEmail}"

        def bodyText = """Build failed.

Branch: ${branch}
Target: ${TARGET_BRANCH}

Job:   ${env.JOB_NAME}
Build: #${env.BUILD_NUMBER}

--- Rebase output ---
${rebaseOutput}

--- Merge output ---
${mergeOutput}

--- Test output ---
${testOutput}

Console: ${env.BUILD_URL}console
"""

        def mailArgs = [
          to: recipients,
          subject: "Jenkins ${currentBuild.currentResult}: ${branch} (target: ${TARGET_BRANCH})",
          body: bodyText
        ]

        if (env.FROM_EMAIL?.trim()) {
          mailArgs.from = env.FROM_EMAIL.trim()
        }

        emailext(mailArgs)
      }
    }

    unstable {
      script {
        // Reuse the failure path so unstable builds also notify admins/developers.
        def branch = (env.EFFECTIVE_BRANCH ?: env.BRANCH_NAME ?: env.GIT_BRANCH ?: "unknown")
        branch = branch.replaceFirst(/^origin\//, "")

        def mergeOutput  = fileExists('merge_output.txt')  ? readFile('merge_output.txt')  : ""
        def rebaseOutput = fileExists('rebase_output.txt') ? readFile('rebase_output.txt') : ""
        def testOutput   = fileExists('test_output.txt')   ? readFile('test_output.txt')   : ""

        def authorEmail = ""
        try {
          authorEmail = sh(script: "git log -1 --pretty=format:%ae || true", returnStdout: true).trim()
        } catch (e) {
          authorEmail = ""
        }

        def recipients = "${ADMIN_EMAIL}"
        if (authorEmail) recipients = "${recipients}, ${authorEmail}"

        def bodyText = """Build result: ${currentBuild.currentResult}.

Branch: ${branch}
Target: ${TARGET_BRANCH}

Job:   ${env.JOB_NAME}
Build: #${env.BUILD_NUMBER}

--- Rebase output ---
${rebaseOutput}

--- Merge output ---
${mergeOutput}

--- Test output ---
${testOutput}

Console: ${env.BUILD_URL}console
"""

        def mailArgs = [
          to: recipients,
          subject: "Jenkins ${currentBuild.currentResult}: ${branch} (target: ${TARGET_BRANCH})",
          body: bodyText
        ]

        if (env.FROM_EMAIL?.trim()) {
          mailArgs.from = env.FROM_EMAIL.trim()
        }

        emailext(mailArgs)
      }
    }
  }
}

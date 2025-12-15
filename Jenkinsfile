pipeline {
  agent any

  environment {
    TARGET_BRANCH      = "main"
    GIT_CREDENTIALS_ID = "github-creds"
    ADMIN_EMAIL        = "rahul.singhh.144@gmail.com"
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
      steps { checkout scm }
    }

    stage('Run tests') {
      when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
      steps {
        // If you add go.mod, switch to: sh 'go test ./...'
        sh '''
          set -e
          if [ -f go.mod ]; then
            go main.go ./...
          else
            echo "go.mod not found; running go test in current package only"
            go main.go ./
          fi
        '''
      }
    }

    stage('Rebase feature onto latest main') {
      when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
      steps {
        script {
          sh """
            git config user.name "jenkins-bot"
            git config user.email "jenkins-bot@example.com"

            git fetch origin

            # Ensure we are on the feature branch locally
            git checkout -B ${env.EFFECTIVE_BRANCH} origin/${env.EFFECTIVE_BRANCH}

            # Rebase feature branch onto latest main
            set +e
            git rebase origin/${TARGET_BRANCH} > rebase_output.txt 2>&1
            REBASE_STATUS=\$?
            set -e

            if [ "\$REBASE_STATUS" -ne 0 ]; then
              echo "Rebase conflict detected!"
              git rebase --abort || true
              exit 98
            fi
          """
        }
      }
    }

    stage('Merge rebased feature into main') {
      when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
      steps {
        script {
          sh """
            git fetch origin

            # Checkout latest main
            git checkout -B ${TARGET_BRANCH} origin/${TARGET_BRANCH}

            # Merge the (rebased) feature branch
            set +e
            git merge --no-ff ${env.EFFECTIVE_BRANCH} > merge_output.txt 2>&1
            MERGE_STATUS=\$?
            set -e

            if [ "\$MERGE_STATUS" -ne 0 ]; then
              echo "Merge conflict detected!"
              git merge --abort || true
              exit 99
            fi

            git log -1 --oneline
          """
        }
      }
    }

    stage('Push main') {
      when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
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
      // Always try to detect branch, but NEVER skip emailing because of it
      def branch = (env.EFFECTIVE_BRANCH ?: env.BRANCH_NAME ?: env.GIT_BRANCH ?: "unknown")
      branch = branch.replaceFirst(/^origin\//, "")

      def mergeOutput  = fileExists('merge_output.txt')  ? readFile('merge_output.txt')  : ""
      def rebaseOutput = fileExists('rebase_output.txt') ? readFile('rebase_output.txt') : ""

      // try to get committer email (may fail if checkout failed)
      def authorEmail = ""
      try {
        authorEmail = sh(script: "git log -1 --pretty=format:%ae || true", returnStdout: true).trim()
      } catch (e) {
        authorEmail = ""
      }

      def recipients = "${ADMIN_EMAIL}"
      if (authorEmail) recipients = "${recipients}, ${authorEmail}"

      def note = branch.startsWith("feat/") ? "" : "\nNOTE: This branch is not feat/* (or branch unknown). Email still sent.\n"

      emailext(
        to: recipients,
        subject: "Jenkins FAILED: ${branch} (target: ${TARGET_BRANCH})",
        body: """Build failed.${note}

Branch: ${branch}
Job: ${env.JOB_NAME}
Build: #${env.BUILD_NUMBER}

--- Rebase output ---
${rebaseOutput}

--- Merge output ---
${mergeOutput}

Console: ${env.BUILD_URL}console
"""
      )
    }
  }
}
}
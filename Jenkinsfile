pipeline {
  agent any

  environment {
    TARGET_BRANCH      = "main"
    // ✅ IMPORTANT: no trailing spaces/tabs in credentials id
    GIT_CREDENTIALS_ID = "6f9dfd84-7cc3-412e-8c73-b6c8bd1a3291"
    ADMIN_EMAIL        = "rahul.singhh.144@gmail.com"
    REPO_URL           = "https://github.com/RahulVervebot/go_cicd_demo.git"
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

    stage('Run tests') {
      when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
      steps {

        sh '''
          set -e
          if [ -f go.mod ]; then
            echo "go.mod found; running: go test ./..."
            go test ./...
          else
            echo "go.mod not found; skipping tests (recommended: add go.mod and enable go test ./...)"
          fi
        '''
      }
    }



stage('Attempt merge into main') {
  when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
  steps {
    sh """
      set -e
      git config user.name "jenkins-bot"
      git config user.email "jenkins-bot@example.com"

      git fetch origin

      # checkout latest main
      git checkout -B ${TARGET_BRANCH} origin/${TARGET_BRANCH}

      # try merge feature
      set +e
      git merge --no-ff origin/${env.EFFECTIVE_BRANCH} > merge_output.txt 2>&1
      MERGE_STATUS=\$?
      set -e

      if [ "\$MERGE_STATUS" -ne 0 ]; then
        echo "Merge conflict detected!"
        git merge --abort || true
        exit 99
      fi
    """
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

            # Askpass helper so secrets aren't injected into the URL
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

            # Keep remote clean (no creds in URL)
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

        // committer email (developer)
        def authorEmail = ""
        try {
          authorEmail = sh(script: "git log -1 --pretty=format:%ae || true", returnStdout: true).trim()
        } catch (e) {
          authorEmail = ""
        }

        def recipients = "${ADMIN_EMAIL}"
        if (authorEmail) recipients = "${recipients}, ${authorEmail}"

        emailext(
          to: recipients,
          subject: "Jenkins SUCCESS: ${branch} merged into ${TARGET_BRANCH}",
          body: """Build & merge successful.

Feature branch: ${branch}
Target branch:  ${TARGET_BRANCH}

Job:   ${env.JOB_NAME}
Build: #${env.BUILD_NUMBER}

Console: ${env.BUILD_URL}console
"""
        )
      }
    }

    failure {
      script {
        def branch = (env.EFFECTIVE_BRANCH ?: env.BRANCH_NAME ?: env.GIT_BRANCH ?: "unknown")
        branch = branch.replaceFirst(/^origin\//, "")

        def mergeOutput  = fileExists('merge_output.txt')  ? readFile('merge_output.txt')  : ""
        def rebaseOutput = fileExists('rebase_output.txt') ? readFile('rebase_output.txt') : ""

        def authorEmail = ""
        try {
          authorEmail = sh(script: "git log -1 --pretty=format:%ae || true", returnStdout: true).trim()
        } catch (e) {
          authorEmail = ""
        }

        def recipients = "${ADMIN_EMAIL}"
        if (authorEmail) recipients = "${recipients}, ${authorEmail}"

        emailext(
          to: recipients,
          subject: "Jenkins FAILED: ${branch} (target: ${TARGET_BRANCH})",
          body: """Build failed.

Branch: ${branch}
Target: ${TARGET_BRANCH}

Job:   ${env.JOB_NAME}
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

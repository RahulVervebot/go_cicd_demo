pipeline {
  agent any

  environment {
    TARGET_BRANCH      = "main"
    GIT_CREDENTIALS_ID = "6f9dfd84-7cc3-412e-8c73-b6c8bd1a3291	"
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
        // If you add go.mod
        sh '''
          set -e
          if [ -f go.mod ]; then
            go test ./...
          else
            echo "go.mod not found; running go test in current package only"
            go test ./
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

stage('Push main') {
  when { expression { (env.EFFECTIVE_BRANCH ?: "").startsWith("feat/") } }
  steps {
    withCredentials([usernamePassword(credentialsId: env.GIT_CREDENTIALS_ID,
      usernameVariable: 'GIT_USER',
      passwordVariable: 'GIT_PASS')]) {

      sh '''
        set -e

        # Create an askpass helper so git can read creds safely (no URL injection)
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

        # Ensure origin is the normal URL (no creds embedded)
        git remote set-url origin https://github.com/RahulVervebot/go_cicd_demo.git

        git push origin main

        rm -f .git_askpass.sh
      '''
    }
  }
}

}
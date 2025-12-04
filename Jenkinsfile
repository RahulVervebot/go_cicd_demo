pipeline {
  agent any

  environment {
    APP_NAME = "go_cicd_demo"
    DOCKER_IMAGE = "go_cicd_demo:latest"
  }

  stages {
    stage('Checkout') {
      steps { checkout scm }
    }

    stage('Go Format Check') {
      steps {
        sh '''
          set -e
          UNFORMATTED=$(gofmt -l .)
          if [ -n "$UNFORMATTED" ]; then
            echo "These files are not gofmt'd:"
            echo "$UNFORMATTED"
            exit 1
          fi
        '''
      }
    }

    stage('Test') {
      steps {
        sh 'go test ./... -v'
      }
    }

    stage('Build') {
      steps {
        sh 'go build -o server ./cmd/api'
      }
    }

    stage('Docker Build') {
      steps {
        sh 'docker build -t ${DOCKER_IMAGE} .'
      }
    }

    stage('Deploy (local demo)') {
      when { branch "main" }
      steps {
        sh '''
          docker rm -f ${APP_NAME} || true
          docker run -d --name ${APP_NAME} -p 8080:8080 ${DOCKER_IMAGE}
        '''
      }
    }
  }

  post {
    always {
      echo "Pipeline finished: ${currentBuild.currentResult}"
    }
  }
}

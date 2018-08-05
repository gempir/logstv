pipeline {
  agent any
  stages {
    stage('build bot') {
      steps {
        sh '''cd bot
go get ./...
go build
'''
      }
    }
  }
}
pipeline {
    agent none
    stages {
        stage('Build') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'cd messages && make regenerate'
                sh 'cd treeservice && go build main.go'
                sh 'cd treecli && go build main.go'
            }
        }
        stage('Build Docker Image') {
            agent any
            steps {
                sh "docker-build-and-push -b ${BRANCH_NAME} -s treeservice -f treeservice.dockerfile"
                sh "docker-build-and-push -b ${BRANCH_NAME} -s treecli -f treecli.dockerfile"
            }
        }
    }
}

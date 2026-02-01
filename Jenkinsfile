pipeline {
    agent { label 'deployer' }

    tools { 
        go 'go1.23'
        'org.jenkinsci.plugins.docker.commons.tools.DockerTool' 'DockerDefault'
        
    }

    environment {
        ALERT_EMAIL = credentials('7f303f7b-bb98-44b0-83b9-333ffe60c73f')
        DOCKER_IMAGE_BASE = 'bluffpointtech.com/bpti18next' // Replace with your Docker image name
        DOCKERHUB_IMAGE_BASE = 'bluffpointtech/bptopenware'
        DOCKERHUB_REGISTRY = 'https://index.docker.io/v2/bluffpointtech/bptopenware' // Docker Hub registry
        LOCAL_REGISTRY = 'https://bptregistry:6721' // Replace with your local Docker registry
        DOCKERHUB_CREDENTIALS = credentials('a0618e87-90fe-4804-807a-2bac6ec201d2') // Replace with your Docker Hub credentials ID
        LOCAL_CREDENTIALS_ID = 'dockerregistry' // Replace with your local Docker credentials ID
    }

    stage('Init Version') {
        steps {
            script {
                def props = readProperties file: '.env'
                env.VERSION = props.VERSION
                env.DOCKER_IMAGE = "${DOCKER_IMAGE_BASE}:${env.VERSION}"
                env.DOCKERHUB_IMAGE = "${DOCKERHUB_IMAGE_BASE}:bpti18next-${env.VERSION}"
            }
        }
    }
    stages {

        stage('SonarQube') {
            steps {
                script { scannerHome = tool 'SonarQube' }
                withSonarQubeEnv('SonarQube') {
                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectVersion=${env.VERSION}"
                }
            }
        }

        stage('Build') {
            steps {
              sh 'make clean && make build'
            }
        }

        stage('Docker Build and push to Local') {
            steps {
                script {
                    docker.withRegistry("${LOCAL_REGISTRY}", "${LOCAL_CREDENTIALS_ID}") {
                        docker.build("${DOCKER_IMAGE}", "--build-arg VERSION=${env.VERSION}").push()
                    }
                }
            }
        }

        stage('docker logout') {
            steps {
              sh 'docker logout'
            }
        }

        stage('Docker Build and push to DockerHub') {
            steps {
                sh "docker tag ${DOCKER_IMAGE} ${DOCKERHUB_IMAGE}"
                sh "docker login -u ${DOCKERHUB_CREDENTIALS_USR} -p ${DOCKERHUB_CREDENTIALS_PSW}"
                sh "docker push ${DOCKERHUB_IMAGE}"
                sh 'docker logout'
            }
        }

    }

    post {
        always {
            // Archive the build artifacts
            archiveArtifacts artifacts: '**/build/**', allowEmptyArchive: true
        }
        success {
            echo 'Build succeeded!'
        }
        failure {
            echo 'Build failed!'
            mail to: "${ALERT_EMAIL}",
                 subject: "Build Failed: ${env.JOB_NAME} [${env.BUILD_NUMBER}]",
                 body: "The build has failed. Please check the Jenkins console output for more details."
        }
    }
}
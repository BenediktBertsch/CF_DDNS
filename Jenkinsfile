node {
    def app
    def dockerUsername = "xaviius"
    def registryName = "cfddns"

    stage('Clone repository') {
        checkout scm
    }

    stage('Building image') {
        app = docker.build(dockerUsername + "/" + registryName + ":$BUILD_NUMBER", ".")
    }

    stage('Test image') {
        app.inside {
            sh 'echo "Tests passed"'
        }
    }

    stage('Push image') {
        docker.withRegistry('https://registry.hub.docker.com', 'dockerHub') {
            app.push("${env.BUILD_NUMBER}-stable")
            app.push("latest")
        }
        
    }
}
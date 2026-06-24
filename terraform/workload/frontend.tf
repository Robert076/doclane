resource "kubernetes_deployment" "frontend" {
  metadata {
    name      = "frontend"
    namespace = kubernetes_namespace.doclane.metadata[0].name
  }

  spec {
    replicas = 2

    selector {
      match_labels = { app = "frontend" }
    }

    template {
      metadata {
        labels = { app = "frontend" }
      }

      spec {
        container {
          name  = "frontend"
          image = "${local.ecr_frontend}:latest"

          port {
            container_port = 3000
          }

          env {
            name  = "BACKEND_URL"
            value = "http://backend.doclane.svc.cluster.local:8080/api"
          }

          env {
            name  = "NEXT_PUBLIC_APP_URL"
            value = "https://thesis.robert-beres.com"
          }

          env {
            name  = "NEXT_PUBLIC_AWS_REGION"
            value = "eu-west-1"
          }

          liveness_probe {
            http_get {
              path = "/api/health"
              port = 3000
            }
            initial_delay_seconds = 10
            period_seconds        = 10
          }

          readiness_probe {
            http_get {
              path = "/api/health"
              port = 3000
            }
            initial_delay_seconds = 5
            period_seconds        = 5
          }

          resources {
            requests = {
              cpu    = "100m"
              memory = "256Mi"
            }
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "frontend" {
  metadata {
    name      = "frontend"
    namespace = kubernetes_namespace.doclane.metadata[0].name
  }

  spec {
    selector = { app = "frontend" }

    port {
      port        = 3000
      target_port = 3000
    }
  }
}

resource "kubernetes_horizontal_pod_autoscaler_v2" "frontend" {
  metadata {
    name      = "frontend"
    namespace = kubernetes_namespace.doclane.metadata[0].name
  }

  spec {
    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = kubernetes_deployment.frontend.metadata[0].name
    }

    min_replicas = 2
    max_replicas = 4

    metric {
      type = "Resource"
      resource {
        name = "cpu"
        target {
          type                = "Utilization"
          average_utilization = 70
        }
      }
    }
  }
}

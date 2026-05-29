# Backend: ServiceAccount (IRSA), ConfigMap, Secret, Deployment, Service, HPA.

resource "kubernetes_service_account" "backend" {
  metadata {
    name      = "doclane-backend"
    namespace = kubernetes_namespace.doclane.metadata[0].name
    annotations = {
      "eks.amazonaws.com/role-arn" = local.backend_pod_role_arn
    }
  }
}

resource "kubernetes_config_map" "backend" {
  metadata {
    name      = "backend-config"
    namespace = kubernetes_namespace.doclane.metadata[0].name
  }

  data = {
    AWS_REGION           = "eu-west-1"
    DB_HOST              = local.rds_address
    DB_PORT              = "5432"
    DB_NAME              = "doclane"
    DB_SSLMODE           = "require"
    S3_BUCKET_NAME       = local.s3_bucket
    COGNITO_USER_POOL_ID = local.cognito_pool_id
    COGNITO_CLIENT_ID    = local.cognito_client_id
    ALLOWED_ORIGIN       = "https://thesis.robert-beres.com"
  }
}

resource "kubernetes_secret" "backend" {
  metadata {
    name      = "backend-secrets"
    namespace = kubernetes_namespace.doclane.metadata[0].name
  }

  data = {
    DB_USER     = var.db_username
    DB_PASSWORD = var.db_password
    SEED_SECRET = var.seed_secret
  }
}

resource "kubernetes_deployment" "backend" {
  metadata {
    name      = "backend"
    namespace = kubernetes_namespace.doclane.metadata[0].name
  }

  spec {
    replicas = 2

    selector {
      match_labels = { app = "backend" }
    }

    template {
      metadata {
        labels = { app = "backend" }
      }

      spec {
        service_account_name = kubernetes_service_account.backend.metadata[0].name

        container {
          name  = "backend"
          image = "${local.ecr_backend}:latest"

          port {
            container_port = 8080
          }

          env_from {
            config_map_ref {
              name = kubernetes_config_map.backend.metadata[0].name
            }
          }

          env_from {
            secret_ref {
              name = kubernetes_secret.backend.metadata[0].name
            }
          }

          liveness_probe {
            http_get {
              path = "/health"
              port = 8080
            }
            initial_delay_seconds = 5
            period_seconds        = 10
          }

          readiness_probe {
            http_get {
              path = "/health"
              port = 8080
            }
            initial_delay_seconds = 3
            period_seconds        = 5
          }

          resources {
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
            limits = {
              cpu    = "500m"
              memory = "256Mi"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "backend" {
  metadata {
    name      = "backend"
    namespace = kubernetes_namespace.doclane.metadata[0].name
  }

  spec {
    selector = { app = "backend" }

    port {
      port        = 8080
      target_port = 8080
    }
  }
}

resource "kubernetes_horizontal_pod_autoscaler_v2" "backend" {
  metadata {
    name      = "backend"
    namespace = kubernetes_namespace.doclane.metadata[0].name
  }

  spec {
    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = kubernetes_deployment.backend.metadata[0].name
    }

    min_replicas = 2
    max_replicas = 6

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

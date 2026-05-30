resource "kubernetes_job" "insert_admin" {
  metadata {
    name      = "insert-admin"
    namespace = kubernetes_namespace.doclane.metadata[0].name
  }

  spec {
    backoff_limit = 5

    template {
      metadata {
        labels = { job = "insert-admin" }
      }

      spec {
        restart_policy = "OnFailure"

        container {
          name  = "insert-admin"
          image = "amazon/aws-cli:latest"

          command = ["sh", "-c", <<-EOT
            echo "Waiting for backend to be ready..."
            until curl -sf http://backend.doclane.svc.cluster.local:8080/health; do
              sleep 2
            done
            echo "Backend ready. Authenticating with Cognito..."
            TOKEN=$(aws cognito-idp initiate-auth \
              --auth-flow USER_PASSWORD_AUTH \
              --client-id "$COGNITO_CLIENT_ID" \
              --auth-parameters USERNAME=admin@admin.com,PASSWORD=Admin1234! \
              --region eu-west-1 \
              --query 'AuthenticationResult.IdToken' \
              --output text)
            if [ -z "$TOKEN" ] || [ "$TOKEN" = "None" ]; then
              echo "Failed to get token from Cognito"
              exit 1
            fi
            echo "Got token. Calling insert-admin..."
            curl -sf -X POST http://backend.doclane.svc.cluster.local:8080/api/auth/insert-admin \
              -H "Authorization: Bearer $TOKEN" \
              -H "X-Seed-Secret: $SEED_SECRET" \
              -H "Content-Type: application/json"
            echo ""
            echo "Done."
          EOT
          ]

          env {
            name  = "SEED_SECRET"
            value = var.seed_secret
          }
          env {
            name  = "COGNITO_CLIENT_ID"
            value = local.cognito_client_id
          }
        }
      }
    }
  }

  wait_for_completion = true
  timeouts {
    create = "5m"
  }

  depends_on = [kubernetes_deployment.backend, kubernetes_job.db_seed]
}

resource "kubernetes_config_map" "init_sql" {
  metadata {
    name      = "init-sql"
    namespace = kubernetes_namespace.doclane.metadata[0].name
  }

  data = {
    "init.sql" = file("${path.module}/../../init.sql")
  }
}

resource "kubernetes_job" "db_seed" {
  metadata {
    name      = "db-seed"
    namespace = kubernetes_namespace.doclane.metadata[0].name
  }

  spec {
    backoff_limit = 3

    template {
      metadata {
        labels = { job = "db-seed" }
      }

      spec {
        restart_policy = "OnFailure"

        container {
          name  = "seed"
          image = "postgres:17-alpine"

          command = ["sh", "-c", "psql -f /sql/init.sql"]

          env {
            name  = "PGHOST"
            value = local.rds_address
          }
          env {
            name  = "PGPORT"
            value = "5432"
          }
          env {
            name  = "PGDATABASE"
            value = "doclane"
          }
          env {
            name  = "PGUSER"
            value = var.db_username
          }
          env {
            name  = "PGPASSWORD"
            value = var.db_password
          }

          volume_mount {
            name       = "sql"
            mount_path = "/sql"
            read_only  = true
          }
        }

        volume {
          name = "sql"
          config_map {
            name = kubernetes_config_map.init_sql.metadata[0].name
          }
        }
      }
    }
  }

  wait_for_completion = true
  timeouts {
    create = "5m"
  }
}

# Namespace for all doclane workloads.

resource "kubernetes_namespace" "doclane" {
  metadata {
    name = "doclane"
  }
}

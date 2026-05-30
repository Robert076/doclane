resource "kubernetes_namespace" "doclane" {
  metadata {
    name = "doclane"
  }
}

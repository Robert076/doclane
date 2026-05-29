# Route53 A record pointing thesis.robert-beres.com at the ALB.

data "aws_lb" "ingress" {
  tags = {
    "elbv2.k8s.aws/cluster"    = local.cluster_name
    "ingress.k8s.aws/resource" = "LoadBalancer"
  }

  depends_on = [kubernetes_ingress_v1.main]
}

resource "aws_route53_record" "app" {
  zone_id = data.terraform_remote_state.data.outputs.route53_zone_id
  name    = "thesis.robert-beres.com"
  type    = "A"

  alias {
    name                   = data.aws_lb.ingress.dns_name
    zone_id                = data.aws_lb.ingress.zone_id
    evaluate_target_health = true
  }
}

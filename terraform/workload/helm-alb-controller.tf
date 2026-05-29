# AWS Load Balancer Controller — watches Ingress objects and creates real ALBs.

resource "helm_release" "alb_controller" {
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  version    = "1.11.0"
  namespace  = "kube-system"

  set {
    name  = "clusterName"
    value = local.cluster_name
  }

  set {
    name  = "serviceAccount.create"
    value = "true"
  }

  set {
    name  = "serviceAccount.name"
    value = "aws-load-balancer-controller"
  }

  set {
    name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
    value = local.alb_controller_role_arn
  }

  set {
    name  = "region"
    value = "eu-west-1"
  }

  set {
    name  = "vpcId"
    value = data.terraform_remote_state.compute.outputs.vpc_id
  }
}

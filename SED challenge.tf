provider "aws" {
  region = "us-east-1" 
}

resource "aws_instance" "web_server" {
  ami           = "ami-0c55b159cbfafe1f0" 
  instance_type = "t2.micro" 
  key_name      = "your-key-pair" 

  tags = {
    Name = "web-server-instance"
  }

}

resource "aws_security_group" "web_server_sg" {
  name        = "web-server-sg"
  description = "Security group for web server"
  
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb" "web_server_lb" {
  name               = "web-server-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.web_server_sg.id]
  subnets            = ["subnet-12345678"] 
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.web_server_lb.arn
  port              = "443"
  protocol          = "HTTPS"
  
  default_action {
    type             = "redirect"
    redirect {
      protocol       = "HTTPS"
      port           = "443"
      status_code    = "HTTP_301"
    }
  }

  ssl_policy = "ELBSecurityPolicy-2016-08"
}

resource "aws_acm_certificate" "web_server_cert" {
  domain_name       = "example.com" 
  validation_method = "DNS"
}

resource "aws_route53_record" "cert_validation" {
  zone_id = "your-zone-id" 
  name    = "_acme-challenge.example.com" 
  type    = "CNAME"
  ttl     = "300"
  records = [aws_acm_certificate.web_server_cert.domain_validation_options.0.resource_record_name]
}

resource "aws_acm_certificate_validation" "web_server_cert_validation" {
  certificate_arn         = aws_acm_certificate.web_server_cert.arn
  validation_record_fqdns = [aws_route53_record.cert_validation.fqdn]
}

/*   */
output "server_ip" {
    value = aws_instance.d99df.*.public_ip
    description = "The public IP address of the main server instance."
}

output "regestry_name" {
    value = aws_ecr_repository.hello-world.repository_url
    description = ""
}
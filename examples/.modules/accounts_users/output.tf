output "idps" {
  value = {
    for idp in data.codefresh_idps.idps:
      idp.id => { client_name = idp.client_name,
                  display_name = idp.display_name 
                }
  }  
}
output "accounts" {
  value = {
    for acc in codefresh_account.acc:
      acc.id => acc.name
  }
}

export function cleanUpJWT(jwtToken: string): string {
	if (jwtToken === "") {
		return ""
	}
	if(jwtToken.startsWith("jwt=")){
    jwtToken = jwtToken.slice(4) //remove "jwt="
  }
	jwtToken = jwtToken.split(' ')[0]
	jwtToken = jwtToken.slice(0, -1) //remove the last ";" character
	return jwtToken
}

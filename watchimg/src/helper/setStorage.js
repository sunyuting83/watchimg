export default (type = true, token) => {
  if (type) {
    localStorage.setItem('token', token)
  }else{
    localStorage.clear()
  }
}
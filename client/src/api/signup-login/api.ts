export const signUp = async (name: string, password: string) => {
    if (name === "" || password == "") {
        throw new Error("need name and password")
    }
    const response = await fetch('http://localhost:8080/signup', {
        method: 'POST',
        body: JSON.stringify({
            name,
            password
        })
    })
    if (!response.ok) {
        throw new Error("Failed to sign up")
    }
    const data = await response.json().catch((err) => {
        console.log("failed to get signup data")
        console.log(response)
        console.log(err)
    })
    return data;
}

export const login = async (name: string, password: string) => {
    if (name === "" || password === "") {
        throw new Error("need name and password")
    }
    const response = await fetch('http://localhost:8080/login', {
        method: 'POST',
        body: JSON.stringify({
            name,
            password
        })  
    })
    if (!response.ok) {
        throw new Error("Failed to login")
    }
    const data = await response.json().catch((err) => {
        console.log("failed to get login data")
        console.log(err)
    })
    return data
}
const headerToken = (token: string) => ({"Authorization": `Bearer : ${token}`})

export const getTodos = async (token: string) => {
    const response = await fetch('http://localhost:8080/todos', {
        headers: headerToken(token)
    });
    if (!response.ok) {
        throw new Error('Failed to fetch todos');
    }
    const data = await response.json();
    return data;
};

export const postTodo = async (token: string, description: string) => {
    const response = await fetch('http://localhost:8080/todos', {
        method: 'POST',
        headers: headerToken(token),
        body: JSON.stringify({
            description: description,
        })
    })
    if (!response.ok) {
        throw new Error(`Failed to post todo`)
    }
}

export const deleteTodo = async (token: string, todo_id: string) => {
    const response = await fetch(`http://localhost:8080/todos/${todo_id}`, {
        method: 'DELETE',
        headers: headerToken(token)
    })
    if (!response.ok) {
        throw new Error(`Failed to delete todo ${todo_id}`)
    }
}

export const editTodo = async (token: string, todo_id: string, description: string) => {
    console.log("todo_id : ", todo_id)
    const response = await fetch(`http://localhost:8080/todos/${todo_id}`, {
        method: 'PUT',
        headers: headerToken(token),
        body: JSON.stringify({
            description: description,
        })
    })
    if (!response.ok) {
        throw new Error(`Failed to update todo`)
    }
}

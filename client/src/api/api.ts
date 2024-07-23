export const getTodos = async () => {
    const response = await fetch('http://localhost:8080/todos');
    if (!response.ok) {
        throw new Error('Failed to fetch todos');
    }
    const data = await response.json();
    return data;
};

export const postTodo = async (description: string) => {
    const response = await fetch('http://localhost:8080/todos', {
        method: 'POST',
        body: JSON.stringify({
            description: description,
        })
    })
    if (!response.ok) {
        throw new Error(`Failed to post todo`)
    }
}

export const deleteTodo = async (todo_id: string) => {
    const response = await fetch(`http://localhost:8080/todos/${todo_id}`, {
        method: 'DELETE',
    })
    if (!response.ok) {
        throw new Error(`Failed to delete todo ${todo_id}`)
    }
}
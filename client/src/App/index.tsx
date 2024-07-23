import { useEffect, useState } from "react";
import "./index.css";
import { deleteTodo, getTodos, postTodo } from "../api/api";

interface Todos {
  todo_id: string;
  description: string;
}

function App() {
  const [todos, setTodos] = useState<Todos[]>();
  const [input, setInput] = useState("");
  const [openModal, setOpenModal] = useState(false)

  const fetchTodos = async () => {
    try {
      const response = await getTodos();
      setTodos(response);
      console.log(typeof response);
    } catch (error) {
      console.log(error);
    }
  };
  useEffect(() => {
    fetchTodos();
  }, []);

  return (
    <main onClick={() => {setOpenModal(false)}}>
      <h1>To-do List</h1>
      <form className="todo-form" onSubmit={async (e) => {
        e.preventDefault()
        await postTodo(input)
        fetchTodos()
        setInput("")
      }}>
        <label htmlFor="description" />
        <input
          type="text"
          id="description"
          name="description"
          placeholder="add your todo here.."
          value={input}
          onChange={(e) => {
            setInput(e.currentTarget.value);
          }}
        />
        <button type="submit" disabled={input === ""}>Add</button>
      </form>
      <div className="todo-list">
        {todos?.map((todo) => (
          <div className="todo" key={todo.todo_id}>
            <p>{todo.description}</p>
            <div className="todo-btns">
              <button onClick={(e) => {
                setOpenModal(true)
                e.stopPropagation()
              }}>
                modifier
              </button>
              <button
                onClick={async () => {
                  await deleteTodo(todo.todo_id);
                  fetchTodos();
                }}
              >
                supprimer
              </button>
            </div>
          </div>
        ))}
         <dialog open={openModal}>
          petit dialog oklm
         </dialog>
      </div>
    </main>
  );
}

export default App;

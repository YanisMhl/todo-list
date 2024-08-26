import { useEffect, useState } from "react";
import "./index.css";
import { deleteTodo, editTodo, getTodos, postTodo } from "../../api/todos/api";
import ReactModal from "react-modal";
import { useNavigate } from "react-router";
interface Todos {
  todo_id: string;
  description: string;
}

interface ModalProps {
  id: string;
  input: string | null;
  opened: boolean;
}

interface UserInfos {
  token: string
  name: string
}

const Todos = () => {
  const navigate = useNavigate()
  const [todos, setTodos] = useState<Todos[]>();
  const [input, setInput] = useState("");
  const userInfos: UserInfos = {token: localStorage.getItem('token') || '', name: localStorage.getItem('name') || ''}
  const [modalState, setModalState] = useState<ModalProps>({
    id: "",
    input: "",
    opened: false,
  });
  const fetchTodos = async () => {
    try {
      const response = await getTodos(userInfos.token);
      setTodos(response);
      console.log(typeof response);
    } catch (error) {
      console.log(error);
    }
  };
  useEffect(() => {
    if (userInfos.name !== null && userInfos.token !== null) {
      fetchTodos()
    }
  }, []);

  const toggleModal = (id?: string) => {
    setModalState({
      id: id ? id : "",
      input: "",
      opened: id !== undefined,
    });
  };

  return (
    <main>
      <h1>To-do List</h1>
      <p onClick={() => {
        localStorage.clear()
        navigate('/register')
      }}>Se déconnecter ?</p>
      <form
        className="todo-form"
        onSubmit={async (e) => {
          e.preventDefault();
          await postTodo(userInfos.token, input);
          fetchTodos();
          setInput("");
        }}
      >
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
        <button type="submit" disabled={input === ""}>
          Add
        </button>
      </form>
      <div className="todo-list">
        {todos?.map((todo) => (
          <div className="todo" key={todo.todo_id}>
            <p>{todo.description}</p>
            <div className="todo-btns">
              <button
                onClick={(e) => {
                  toggleModal(todo.todo_id);
                  e.stopPropagation();
                }}
              >
                modify
              </button>
              <button
                onClick={async () => {
                  await deleteTodo(userInfos.token, todo.todo_id);
                  fetchTodos();
                }}
              >
                delete
              </button>
            </div>
            <ReactModal
              className="modal"
              isOpen={modalState.opened}
              shouldCloseOnEsc
              shouldReturnFocusAfterClose
              onRequestClose={() => {
                toggleModal();
              }}
            >
              <form
                onSubmit={async (e) => {
                  e.preventDefault();
                  if (modalState.input) {
                    console.log("modal state avant put : ");
                    console.log(modalState);
                    await editTodo(userInfos.token, modalState.id, modalState.input);
                  }
                  fetchTodos();
                  toggleModal();
                  setModalState((prev) => {
                    return {
                      ...prev,
                      input: "",
                    };
                  });
                }}
              >
                <input
                  type="text"
                  id="edit"
                  name="edit"
                  placeholder="modify your todo here.."
                  value={modalState.input || ""}
                  onChange={(e) => {
                    const input = e.currentTarget.value;
                    setModalState((prev) => {
                      return {
                        ...prev,
                        input,
                      };
                    });
                  }}
                />
                <button type="submit" disabled={modalState.input === ""}>
                  Update
                </button>
              </form>
            </ReactModal>
          </div>
        ))}
      </div>
    </main>
  );
};

export default Todos;
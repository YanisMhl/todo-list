import { createBrowserRouter } from "react-router-dom";
import Todos from "./Todos";
import Register from "./Register";

const router = createBrowserRouter([
    {
      path: '/todos',
      element: <Todos />, 
    },
    {
      path: '/register',
      element: <Register />
    },
    {
      path: '*',
      element: <Register />
    }
  ])

export default router
import { useNavigate, useParams } from "react-router-dom"
import { API_URL } from "./App"

export const ConfirmationPage = ()=>{
    const {token = ''} = useParams()
    const redirect = useNavigate()
    const handleClick = async () => {
       const response = await fetch(`${API_URL}/users/activate/${token}`, {
           method: 'PUT'
       })

       if (response.ok) {
           redirect('/')
       } else{
           console.error(response)
       }
    }
    return (
       <div>
        <h1>Confirmation</h1>
        <button onClick={handleClick}>Click to confirm</button>
       </div>
    )
}
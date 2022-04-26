import { useState } from 'react';
import axios from "axios";
import { Main } from '@/templates/Main';
import Loading from '@/components/Loading/Loading';

const Index = () => {
  // const [show, setShow] = useState<boolean>(false);
  const [tty, setTty] = useState('');
  const [load, setLoad] = useState(false);
  const BACKEND_URL = process.env.BACKEND_URL || 'https://k8tty.yad2.io';
  const WETTY_URL = process.env.WETTY_URL || "https://k8tty.yad2.io/";

  const update = () => {
    console.log("create")
    setLoad(true)
    if(tty){
      axios
          .post(BACKEND_URL + "/delete", {uid: tty})
          .then(() => {
            setTty('');
            setLoad(false)
          })
          .catch((err) => {
            console.log(err)
            setTty('');
            setLoad(false)
          })
      return
    }
    axios
        .get(BACKEND_URL + "/create",{headers: {
          'Content-Type': 'application/json'
      }})
        .then((res) => {
          setTty(res.data.uid);
          setLoad(false)
        })
        .catch((err) => {
          console.log(err)
          setLoad(false)
        })
  };

  return (
    <Main>
        <div className='w-full h-full m-auto block flex  justify-center'>
          <div className='w-full h-full'>
            <div className='flex justify-center'>
              <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full m-4" onClick={update} disabled={load}>
                {tty ? "delete me" : "Show me"}
              </button>
            </div>
            {load ? <Loading></Loading>: null}
            {tty ? <iframe className='w-[90%] h-[90%] m-auto' src={WETTY_URL + tty + "/"} title="W3Schools Free Online Web Tutorials"></iframe> : null}
            {tty ? <p>{tty}</p> : null}
          </div>
        </div>
    </Main>
  );
};

export default Index;

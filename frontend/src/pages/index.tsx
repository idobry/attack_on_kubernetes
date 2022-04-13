import { useState } from 'react';
import axios from "axios";
import { Main } from '@/templates/Main';

const Index = () => {
  // const [show, setShow] = useState<boolean>(false);
  const [tty, setTty] = useState('');
  const BACKEND_URL = process.env.BACKEND_URL || 'https://k8tty.yad2.io';

  const update = () => {
    console.log("create")
    if(tty){
      setTty('');
      return
    }
      axios
          .get(BACKEND_URL + "/create")
          .then((res) => {
            setTty(res.data.uid);
          });
  };

  return (
    <Main>
        <div className='w-full h-full m-auto block flex  justify-center'>
          <div className='w-full h-full'>
            <div className='flex justify-center'>
              <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full m-4" onClick={update}>
                {tty ? "delete me" : "Show me"}
              </button>
            </div>
            {tty ? <iframe className='w-[90%] h-[90%] m-auto' src={process.env.WETTY_URL + tty} title="W3Schools Free Online Web Tutorials"></iframe> : null}
            {tty ? <p>{tty}</p> : null}
          </div>
        </div>
    </Main>
  );
};

export default Index;

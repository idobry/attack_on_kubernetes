import { useRouter } from 'next/router';
import { useState } from 'react';
import { Meta } from '@/layout/Meta';
import { Main } from '@/templates/Main';

const Index = () => {
  const router = useRouter();
  const [show, setShow] = useState<boolean>(false);
  const handleClick = () => {
    setShow(prevCheck => !prevCheck)
  }

  return (
    <Main>
        <div className='w-full h-full m-auto block flex  justify-center'>
          <div className='w-full h-full'>
            <div className='flex justify-center'>
              <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full m-4" onClick={handleClick}>
                  show me
              </button>
            </div>
            {show ? <iframe className='w-[90%] h-[90%] m-auto' src="http://localhost:3001/wetty" title="W3Schools Free Online Web Tutorials"></iframe> : null}
            
          </div>
        </div>
    </Main>
  );
};

export default Index;

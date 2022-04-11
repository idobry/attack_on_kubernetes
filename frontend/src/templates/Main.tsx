import { ReactNode } from 'react';
import { useRouter } from 'next/router';
import Link from 'next/link';

import { AppConfig } from '@/utils/AppConfig';
import NavBar from '@/components/NavBar/NavBar';

type IMainProps = {
  meta?: ReactNode;
  children?: ReactNode;
};

const Main = (props: IMainProps) => {
  const router = useRouter();
  return (
  <div className="w-screen h-screen text-gray-700 antialiased">
    <NavBar pathName={router.pathname}></NavBar>
    <div className="w-full h-full p-4">
      <div className="w-full h-full rounded bg-gray-800">
        {props.children}
      </div>
    </div>
  </div>
  );
};

export { Main };

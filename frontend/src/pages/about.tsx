import { Meta } from "@/layout/Meta";
import { Main } from "@/templates/Main";

const About = () => (
  <Main meta={<Meta title="Lorem ipsum" description="Lorem ipsum" />}>
      <div className="p-10 mb-8">
        <p className="text-5xl text-white font-bold">How to open terminal</p>
        <p className="text-3xl text-gray-200">
          Click the button bro
        </p>
      </div>
  </Main>
);

export default About;

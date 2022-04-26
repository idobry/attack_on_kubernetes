import Lottie from 'lottie-react';
import animationData from './cat.json';

const Loading = () => {
    const style = {
        height: 400,
        width: 400,
        margin: 'auto'
      };
    const defaultOptions = {
        loop: true,
        autoplay: true,
        animationData: animationData,
        rendererSettings: {
          preserveAspectRatio: "xMidYMid slice"
        }
      };
    return (
        <div className='w-[90%] h-[90%] m-auto'>
            <Lottie 
                animationData={animationData}
                style={style}
                loop={true}
            />
        </div>
    );
};
export default Loading;
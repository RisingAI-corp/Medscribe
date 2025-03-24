import { useRef } from 'react';
import Autoplay from 'embla-carousel-autoplay';
import { Carousel } from '@mantine/carousel';
import '@mantine/carousel/styles.css';

function LandingCarousel() {
  // Configure autoplay for reliable looping
  const autoplay = useRef(Autoplay({ 
    delay: 2000,
    stopOnInteraction: false,
    stopOnMouseEnter: true
  }));
  
  return (
    <Carousel
      height={100}
      slideSize="20%"
      slideGap="xl"
      align="start"
      slidesToScroll={1}
      withControls={false}
      loop
      draggable={false}
      plugins={[autoplay.current]}
      style={{ width: '100%' }}
    >
      {/* Using text slides instead of images */}
      <Carousel.Slide>
        <div className="flex items-center justify-center h-full">
          <span className="text-lg font-bold text-blue-600 bg-blue-100 px-4 py-2 rounded-full">In Beta</span>
        </div>
      </Carousel.Slide>
      <Carousel.Slide>
        <div className="flex items-center justify-center h-full">
          <span className="text-lg font-bold text-green-600 bg-green-100 px-4 py-2 rounded-full">In Beta</span>
        </div>
      </Carousel.Slide>
      <Carousel.Slide>
        <div className="flex items-center justify-center h-full">
          <span className="text-lg font-bold text-purple-600 bg-purple-100 px-4 py-2 rounded-full">In Beta</span>
        </div>
      </Carousel.Slide>
      <Carousel.Slide>
        <div className="flex items-center justify-center h-full">
          <span className="text-lg font-bold text-amber-600 bg-amber-100 px-4 py-2 rounded-full">In Beta</span>
        </div>
      </Carousel.Slide>
      <Carousel.Slide>
        <div className="flex items-center justify-center h-full">
          <span className="text-lg font-bold text-red-600 bg-red-100 px-4 py-2 rounded-full">In Beta</span>
        </div>
      </Carousel.Slide>
      {/* Duplicate slides to ensure smooth looping */}
      <Carousel.Slide>
        <div className="flex items-center justify-center h-full">
          <span className="text-lg font-bold text-blue-600 bg-blue-100 px-4 py-2 rounded-full">In Beta</span>
        </div>
      </Carousel.Slide>
      <Carousel.Slide>
        <div className="flex items-center justify-center h-full">
          <span className="text-lg font-bold text-green-600 bg-green-100 px-4 py-2 rounded-full">In Beta</span>
        </div>
      </Carousel.Slide>
    </Carousel>
  );
}

export default LandingCarousel;
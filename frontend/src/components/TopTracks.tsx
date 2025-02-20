import { Alert, Carousel, Skeleton } from "antd";
import { useQuery } from "@tanstack/react-query";
import { CarouselRef } from "antd/es/carousel";
import { useRef } from "react";
import { CarouselArrow } from "./common/CarouselArrow";
import Title from "antd/es/typography/Title";
import { getTopTracks } from "../services/trackService";
import { TrackCard } from "./common/TrackCard";

const carouselSettings = {
  slidesToShow: 5,
  slidesToScroll: 1,
  dots: false,
  responsive: [
    {
      breakpoint: 1024,
      settings: {
        slidesToShow: 4,
      },
    },
    {
      breakpoint: 768,
      settings: {
        slidesToShow: 2,
      },
    },
    {
      breakpoint: 480,
      settings: {
        slidesToShow: 1,
      },
    },
  ],
};

export const TopTracks = () => {
  const carouselRef = useRef<CarouselRef>(null);
  const {
    data: tracks,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["topTracks"],
    queryFn: getTopTracks,
  });

  if (isLoading) {
    return (
      <div style={{ textAlign: "center", padding: "20px" }}>
        <Skeleton active />
      </div>
    );
  }

  if (isError) {
    return (
      <Alert
        type="error"
        message="Error"
        description="Failed to load top artists. Please try again later."
        showIcon
      />
    );
  }

  if (!tracks || tracks.length === 0) {
    return (
      <Alert
        type="info"
        message="No top tracks"
        description="You haven't listened to any tracks yet."
        showIcon
      />
    );
  }

  return (
    <>
      <Title
        style={{
          textAlign: "center",
          marginBottom: "24px",
        }}
        level={3}
      >
        Here's what you've been listening to
      </Title>
      <div style={{ position: "relative", padding: "0 20px" }}>
        <CarouselArrow
          type="prev"
          onClick={() => carouselRef.current?.prev()}
        />
        <Carousel ref={carouselRef} {...carouselSettings}>
          {tracks?.map((track) => (
            <TrackCard key={track.id} track={track} />
          ))}
        </Carousel>
        <CarouselArrow
          type="next"
          onClick={() => carouselRef.current?.next()}
        />
      </div>
    </>
  );
};

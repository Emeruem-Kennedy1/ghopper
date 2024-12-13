import { Alert, Carousel, Skeleton } from "antd";
import { useQuery } from "@tanstack/react-query";
import { CarouselRef } from "antd/es/carousel";
import { useRef } from "react";
import { getTopArtists } from "../services/artistService";
import { CarouselArrow } from "./common/CarouselArrow";
import { ArtistCard } from "../components/common/ArtistCard";
import Title from "antd/es/typography/Title";

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

export const TopArtists = () => {
  const carouselRef = useRef<CarouselRef>(null);
  const {
    data: artists,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["topArtists"],
    queryFn: getTopArtists,
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

  return (
    <>
      <Title
        style={{
          textAlign: "center",
          marginBottom: "24px",
        }}
        level={3}
      >
        Your Top Artists
      </Title>
      <div style={{ position: "relative", padding: "0 20px" }}>
        <CarouselArrow
          type="prev"
          onClick={() => carouselRef.current?.prev()}
        />
        <Carousel ref={carouselRef} {...carouselSettings}>
          {artists?.map((artist) => (
            <ArtistCard key={artist.id} artist={artist} />
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

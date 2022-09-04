import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  Spinner,
} from "@chakra-ui/react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { BACKEND_URL } from "../utils/constants";

export const Form = () => {
  interface FormData {
    youtubeURL: string;
  }
  const { mutate, isRequesting, isError } = useDownload();
  const {
    register,
    handleSubmit,
    formState: {},
  } = useForm<FormData>();

  return (
    <Box>
      <form onSubmit={handleSubmit((data) => mutate(data.youtubeURL))}>
        <FormControl>
          <FormLabel>Youtube URL</FormLabel>
          <Input
            type="text"
            {...register("youtubeURL", {
              pattern:
                /^(?:https?:\/\/)?(?:www\.)?(?:youtu\.be\/|youtube\.com\/(?:embed\/|v\/|watch\?v=|watch\?.+&v=))((\w|-){11})(?:\S+)?$/,
              required: true,
            })}
          />
          {isRequesting ? (
            <Spinner size={"xl"} />
          ) : (
            <Button type="submit" disabled={isRequesting}>
              Submit
            </Button>
          )}
        </FormControl>
      </form>
    </Box>
  );
};

const useDownload = () => {
  const [isRequesting, setIsRequesting] = useState(false);
  const [isError, setIsError] = useState(false);

  interface ServerResponse {
    success: boolean;
    message: string;
    location?: string;
  }

  const mutate = (videoURL: string) => {
    const startRequest = async () => {
      setIsRequesting(true);
      setIsError(false);
      try {
        const response = await fetch(`${BACKEND_URL}/convert`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            youtube_url: videoURL,
          }),
        });
        if (response.ok) {
          const data: ServerResponse = await response.json();
          if (data && data.success) {
            window.open(data.location, "_blank")?.focus();
            window.location.reload();
          }
        } else {
          setIsError(true);
        }
      } catch (error) {
        setIsError(true);
      }
    };
    startRequest().then(() => {
      setIsRequesting(false);
    });
  };
  return { isRequesting, isError, mutate };
};

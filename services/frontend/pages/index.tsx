import { Box, Center, Heading } from "@chakra-ui/react";
import type { NextPage } from "next";
import { Form } from "../components/Form";
import { Navbar } from "../components/Navbar";

const Home: NextPage = () => {
  return (
    <Box>
      <Navbar />
      <Center>
        <Heading pt="25">
          <Form />
        </Heading>
      </Center>
    </Box>
  );
};

export default Home;

import { Box, Center, Flex, Text } from "@mantine/core";
import TopBar from "../common/TopBar";

export default function None() {
  return (
    <Center w="100vw">
      <Flex direction="column">
        <TopBar />
        <Box>
          <Text>Select something...</Text>
        </Box>
      </Flex>
    </Center>
  );
}

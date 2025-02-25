import {
  Textarea,
  Accordion,
  CopyButton,
  Button,
  LoadingOverlay,
  Box,
} from '@mantine/core';
import { IconDeviceFloppy, IconDownload } from '@tabler/icons-react';
import { useRef, useEffect, useState } from 'react';
import { useDisclosure } from '@mantine/hooks';

interface SoapEditableSectionProps {
  title: string;
  text: string;
  isLoading: boolean;
  handleSave: (newText: string) => void;
  handleLearnFormat: (contentSection: string, content: string) => void;
  handleAutoSave: (newText: string) => void;
}

const SoapSectionBox = ({
  title,
  text,
  isLoading,
  handleSave,
  handleLearnFormat,
  handleAutoSave,
}: SoapEditableSectionProps) => {
  const [clicked, setClicked] = useState(false);

  const [visible] = useDisclosure(false);
  const [isDirty, setIsDirty] = useState(false);

  const latestLocalTextRef = useRef(text);

  useEffect(() => {
    return () => {
      if (latestLocalTextRef.current !== text) {
        handleAutoSave(latestLocalTextRef.current);
      }
    };
  }, []);

  const handleClick = () => {
    console.log('learning style');
    handleLearnFormat(title, latestLocalTextRef.current);
    setClicked(true);
    setTimeout(() => {
      setClicked(false);
    }, 2000);
  };

  return (
    <div className="relative">
      <Accordion
        defaultValue={title}
        className="shadow-lg rounded-lg overflow-hidden bg-white"
      >
        <Accordion.Item key={title} value={title}>
          <Accordion.Control>
            <div className="flex justify-between items-center">
              <span className="text-xl font-bold">{title}</span>
            </div>
          </Accordion.Control>
          <Accordion.Panel>
            <Textarea
              key={text}
              defaultValue={text}
              autosize={true}
              maxRows={10}
              minRows={4}
              placeholder="Input Text Here"
              onChange={e => {
                if (e.currentTarget.value !== text) {
                  setIsDirty(true);
                } else {
                  setIsDirty(false);
                }
                latestLocalTextRef.current = e.currentTarget.value;
              }}
            />

            <LoadingOverlay
              visible={isLoading}
              zIndex={1000}
              overlayProps={{ radius: 'sm', blur: 2 }}
              loaderProps={{ color: 'blue', type: 'bars' }}
            />

            <div className="flex justify-between mt-4">
              <Button
                variant="outline"
                color={clicked ? 'teal' : undefined}
                onClick={handleClick}
              >
                Learn Style
              </Button>

              <div className="flex gap-4">
                <Box className="relative h-full w-6">
                  <LoadingOverlay
                    visible={visible}
                    zIndex={1000}
                    overlayProps={{ radius: 'sm', blur: 2 }}
                    loaderProps={{ color: 'blue', type: 'bars', size: 'sm' }}
                  />
                </Box>
                {isDirty && (
                  <Button
                    variant="outline"
                    rightSection={<IconDeviceFloppy size={20} />}
                    onClick={() => {
                      handleSave(latestLocalTextRef.current);
                      setIsDirty(false);
                    }}
                  >
                    Save
                  </Button>
                )}
                {text !== '' && (
                  <CopyButton value={text}>
                    {({ copied, copy }) => (
                      <Button
                        color={copied ? 'teal' : 'blue'}
                        onClick={copy}
                        rightSection={<IconDownload size={20} />}
                      >
                        {copied ? 'Copied' : 'Copy'}
                      </Button>
                    )}
                  </CopyButton>
                )}
              </div>
            </div>
          </Accordion.Panel>
        </Accordion.Item>
      </Accordion>
    </div>
  );
};

export default SoapSectionBox;

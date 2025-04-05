import {
  Textarea,
  Accordion,
  CopyButton,
  Button,
  LoadingOverlay,
  Loader,
} from '@mantine/core';

import { IconDeviceFloppy, IconDownload } from '@tabler/icons-react';
import { useRef, useEffect, useState } from 'react';

interface SoapEditableSectionProps {
  reportID: string;
  title: string;
  text: string;
  isLoading: boolean;
  handleSave?: (field: string, newText: string, reportID: string) => void;
  handleLearnFormat?: (
    contentSection: string,
    previous: string,
    content: string,
  ) => void;
  onExpand?: () => void;
  isExpanded: boolean;
  readonly: boolean;
  sectionType?: string;
  isLearnStyleLoading?: boolean;
  isContentSaveLoading?: boolean;
}

const SoapSectionBox = ({
  reportID,
  title,
  text,
  isLoading,
  handleSave,
  handleLearnFormat,
  onExpand,
  isExpanded,
  readonly,
  sectionType,
  isContentSaveLoading,
  isLearnStyleLoading,
}: SoapEditableSectionProps) => {
  console.log('saving loading ', isContentSaveLoading);
  const [clicked, setClicked] = useState(false);

  const [isDirty, setIsDirty] = useState(false);

  const latestLocalTextRef = useRef(text);

  useEffect(() => {
    return () => {
      if (latestLocalTextRef.current !== text && sectionType && handleSave) {
        handleSave(sectionType, latestLocalTextRef.current, reportID);
      }
    };
  }, []);

  const handleClick = () => {
    if (handleLearnFormat && sectionType) {
      handleLearnFormat(sectionType, text, latestLocalTextRef.current);
    }
    setClicked(true);
    setTimeout(() => {
      setClicked(false);
    }, 2000);
  };

  return (
    <div className="relative">
      <Accordion
        defaultValue={isExpanded || text ? title : ''}
        className="shadow-lg rounded-lg overflow-hidden bg-white"
        onChange={value => {
          console.log('Accordion value:', value);
          if (text === '' && onExpand && value === title) {
            onExpand();
          }
        }}
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
              maxRows={18}
              minRows={0}
              placeholder="Input Text Here"
              onChange={e => {
                if (e.currentTarget.value !== text) {
                  setIsDirty(true);
                } else {
                  setIsDirty(false);
                }
                latestLocalTextRef.current = e.currentTarget.value;
              }}
              readOnly={readonly}
            />

            <LoadingOverlay
              visible={isLoading}
              zIndex={1000}
              overlayProps={{ radius: 'sm', blur: 2 }}
              loaderProps={{ color: 'blue', type: 'bars' }}
            />

            <div className="flex justify-between items-center mt-4">
              {!readonly && (
                <div className="flex gap-4">
                  {latestLocalTextRef.current !== '' && (
                    <>
                      <div className="flex gap-2 justify-center items-center">
                        <Button
                          variant="outline"
                          color={clicked ? 'teal' : undefined}
                          onClick={handleClick}
                        >
                          Learn Style
                        </Button>
                        {isLearnStyleLoading && (
                          <Loader color={'blue'} size={18} />
                        )}
                      </div>
                    </>
                  )}

                  {isDirty && (
                    <div className="flex gap-2 justify-center items-center">
                      <Button
                        variant="outline"
                        rightSection={<IconDeviceFloppy size={20} />}
                        onClick={() => {
                          if (handleSave && sectionType) {
                            handleSave(
                              sectionType,
                              latestLocalTextRef.current,
                              reportID,
                            );
                            setIsDirty(false);
                          }
                        }}
                      >
                        Save
                      </Button>
                      {isContentSaveLoading && (
                        <Loader color={'blue'} size={18} />
                      )}
                    </div>
                  )}
                </div>
              )}
              <div>
                {latestLocalTextRef.current !== '' && (
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

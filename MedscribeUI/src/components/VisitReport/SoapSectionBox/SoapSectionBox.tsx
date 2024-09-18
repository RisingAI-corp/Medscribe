import { Textarea, Accordion, CopyButton, Button, LoadingOverlay, Box } from '@mantine/core';
import { IconDownload } from '@tabler/icons-react';

interface SoapEditableSectionProps {
  title: string;
  text: string;
  isSaving: boolean;
  timestamp?: string;
  handleSave: (newText: string) => void;
  handleLearnFormat: () => void;
}

const SoapSectionBox = ({
  title,
  text,
  isSaving,
  timestamp = '',
  handleSave,
  handleLearnFormat,
}: SoapEditableSectionProps) => {
  return (
    <>
      <Accordion
        style={{
          boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1), 0 1px 3px rgba(0, 0, 0, 0.08)',
          borderRadius: '8px',
          overflow: 'hidden',
        }}
      >
        <Accordion.Item key={title} value={title}>
          <Accordion.Control>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span style={{ fontSize: '1.25rem', fontWeight: 'bold' }}>{title}</span>
              <span style={{ fontSize: '1rem', paddingRight: '12px' }}>{timestamp}</span>
            </div>
          </Accordion.Control>
          <Accordion.Panel>
            <Textarea
              value={text}
              autosize={true}
              maxRows={10}
              minRows={4}
              placeholder="Input Text Here"
              onChange={event => handleSave(event.currentTarget.value)}
            />

            <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '1rem' }}>
              <Button variant="outline" onClick={handleLearnFormat}>
                {'Learn Format'}
              </Button>

              <div style={{ display: 'flex', gap: '1rem' }}>
                <Box pos="relative" h="100%" w="25">
                  <LoadingOverlay
                    visible={isSaving}
                    zIndex={1000}
                    overlayProps={{ radius: 'sm', blur: 2 }}
                    loaderProps={{ color: 'blue', type: 'bars', size: 'sm' }}
                  />
                </Box>
                <CopyButton value={text}>
                  {({ copied, copy }) => (
                    <Button color={copied ? 'teal' : 'blue'} onClick={copy} rightSection={<IconDownload size={20} />}>
                      {copied ? 'Copied' : 'Copy'}
                    </Button>
                  )}
                </CopyButton>
              </div>
            </div>
          </Accordion.Panel>
        </Accordion.Item>
      </Accordion>
    </>
  );
};

export default SoapSectionBox;

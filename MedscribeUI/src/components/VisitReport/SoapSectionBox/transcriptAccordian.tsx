// TranscriptAccordion.tsx
import { Accordion, Text, ScrollArea, Box } from '@mantine/core';
import { DiarizedTranscript } from '../../../api/serverResponseTypes';

interface TranscriptAccordionProps {
  reportID: string | number;
  title: string;
  transcriptTurns: DiarizedTranscript | undefined | null;
  isLoading?: boolean;
  onExpand?: () => void;
  isExpanded?: boolean;
  readonly?: boolean;
}

const TranscriptAccordion = ({
  title,
  transcriptTurns,
  isLoading,
  onExpand,
  isExpanded,
}: TranscriptAccordionProps) => {
  console.log('accordain expanded');

  return (
    <div className="relative">
      <Accordion
        defaultValue={
          isExpanded || (transcriptTurns && transcriptTurns.length > 0)
            ? title
            : ''
        }
        className="shadow-lg rounded-lg overflow-hidden"
        style={{ backgroundColor: '#f5f5f5' }} // Add background color here
        onChange={value => {
          if (
            (transcriptTurns === null || transcriptTurns?.length === 0) &&
            onExpand &&
            value === title
          ) {
            onExpand();
          }
        }}
      >
        <Accordion.Item key={title} value={title}>
          <Accordion.Control>
            <div className="flex justify-between items-center">
              <span className="text-xl font-bold">
                {title} {isLoading && '(Loading...)'}
              </span>
            </div>
          </Accordion.Control>
          <Accordion.Panel>
            {transcriptTurns && transcriptTurns.length > 0 ? (
              <ScrollArea
                style={{
                  maxHeight: 300,
                  border: '1px solid #ccc',
                  padding: 10,
                  backgroundColor: '#f5f5f5', // Add background color here as well for the scrollable area
                }}
              >
                {transcriptTurns.map((turn, index) => (
                  <Box
                    key={index}
                    mb={5}
                    style={{ backgroundColor: '#f5f5f5' }}
                  >
                    {' '}
                    {/* Add background color here for each turn */}
                    <Text
                      inline
                      c={
                        turn.speaker.toLowerCase() === 'provider'
                          ? 'blue'
                          : turn.speaker.toLowerCase() === 'patient'
                            ? 'red'
                            : 'gray'
                      }
                      fw="bold"
                    >
                      {turn.speaker}
                      {/* : [{formatTime(turn.startTime)} -{' '}
                      {formatTime(turn.endTime)}] */}
                    </Text>
                    <Text inline></Text> <Text inline>{turn.text}</Text>
                    <br />
                  </Box>
                ))}
              </ScrollArea>
            ) : (
              <Text c="dimmed">No transcript available.</Text>
            )}
          </Accordion.Panel>
        </Accordion.Item>
      </Accordion>
    </div>
  );
};

export default TranscriptAccordion;

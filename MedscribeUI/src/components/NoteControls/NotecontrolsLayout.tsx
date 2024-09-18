import { Select, Button } from '@mantine/core';
import { useState, useEffect } from 'react';
import BtnGroupSelector from '../Utilities/BtnGroupSelector';
import { useDisclosure } from '@mantine/hooks';
import { LoadingOverlay } from '@mantine/core';

const defaultVisitType = 'New Patient';
const defaultPronoun = 'HE';
const defaultPatientClient = 'Patient';

function NoteControlsLayout({ isStatus }: { isStatus: boolean }) {
  const [selectedPronoun, setSelectedPronoun] = useState(defaultPronoun);
  const [selectedVisitType, setSelectedVisitType] = useState(defaultVisitType);
  const [selectedPatientClient, setSelectedPatientClient] = useState(defaultPatientClient);
  const [isDirty, setIsDirty] = useState(false);
  const [visible, { toggle }] = useDisclosure(isStatus);

  useEffect(() => {
    if (
      selectedVisitType !== defaultVisitType ||
      selectedPronoun !== defaultPronoun ||
      selectedPatientClient !== defaultPatientClient
    ) {
      setIsDirty(true);
    } else {
      setIsDirty(false);
    }
  }, [selectedVisitType, selectedPronoun, selectedPatientClient]);

  const handleVisitTypeSelect = (value: string | null) => {
    if (value === 'New Patient' || value === 'Returning Patient') {
      setSelectedVisitType(value);
      return;
    }
    setSelectedVisitType(defaultVisitType);
  };

  const handlePronounSelect = (value: string) => {
    if (value === 'SHE' || value === 'HE' || value === 'THEY') {
      setSelectedPronoun(value);
      return;
    }
    setSelectedPronoun(defaultPronoun);
  };

  const handlePatientClientSelect = (value: string) => {
    if (value === 'Patient' || value === 'Client') {
      setSelectedPatientClient(value);
      return;
    }
    setSelectedPatientClient(defaultPatientClient);
  };

  const handleRegenerate = () => {
    toggle();
  };

  return (
    <>
      <LoadingOverlay
        visible={visible}
        zIndex={1000}
        overlayProps={{ radius: 'sm', blur: 2 }}
        loaderProps={{ color: 'blue', type: 'bars' }}
      />

      <span style={{ display: 'block', marginBottom: '8px' }}>Visit Type</span>
      <Select
        defaultValue={selectedVisitType}
        data={['New Patient', 'Returning Patient']}
        value={selectedVisitType}
        onChange={handleVisitTypeSelect}
      />

      <hr style={{ margin: '16px 0' }} />

      <span style={{ display: 'block', marginBottom: '8px' }}>Pronoun Selector</span>
      <BtnGroupSelector
        buttonLabelOptions={['HE', 'SHE', 'THEY']}
        selectedBtn={selectedPronoun}
        onSelect={handlePronounSelect}
      />

      <span style={{ display: 'block', margin: '16px 0 8px' }}>Patient/Client</span>
      <BtnGroupSelector
        buttonLabelOptions={['Patient', 'Client']}
        selectedBtn={selectedPatientClient}
        onSelect={handlePatientClientSelect}
      />

      <hr style={{ margin: '16px 0' }} />

      <Button onClick={handleRegenerate} style={{ width: '100%' }} disabled={!isDirty}>
        Regenerate Report
      </Button>
    </>
  );
}

export default NoteControlsLayout;

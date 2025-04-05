import ProfileSummaryCard from '../ProfileSummaryCard/profileSummaryCard';
import VisitReportLayout from '../VisitReport/visitReportLayout';
import NoteControlsLayout from '../NoteControls/noteControlsLayout';
import { SelectedPatientHeaderInformationAtom } from './derivedAtoms';
import { UpdateSelectedPatientNameAtom } from '../../states/patientsAtom';
import { useAtom } from 'jotai';
import { useMutation } from '@tanstack/react-query';
import { changeName } from '../../api/changeName';

function PatientDashBoardLayout() {
  const [headerInformation, _] = useAtom(SelectedPatientHeaderInformationAtom);
  const [, updateHeaderInformation] = useAtom(UpdateSelectedPatientNameAtom);
  const updateNameMutation = useMutation({
    mutationFn: changeName,
    onError: error => {
      console.error(error);
    },
  });

  const handleUpdateName = (name: string) => {
    updateNameMutation.mutate({
      ReportID: headerInformation.id,
      NewName: name,
    });
  };

  return (
    <div className="flex flex-col h-full max-h-screen overflow-hidden">
      <div className="flex-shrink-0 border-b border-gray-300">
        <ProfileSummaryCard
          name={headerInformation.name}
          description={headerInformation.condensedSummary}
          onChange={updateHeaderInformation}
          handleUpdateName={handleUpdateName}
        />
      </div>

      <div className="flex-1 overflow-y-auto p-4">
        <div className="flex gap-16">
          <div className="flex-1 pb-24">
            <VisitReportLayout />
          </div>

          <div className="flex-2 pr-3 flex-shrink-0 min-w-[300px]">
            <NoteControlsLayout />
          </div>
        </div>
      </div>
    </div>
  );
}

export default PatientDashBoardLayout;

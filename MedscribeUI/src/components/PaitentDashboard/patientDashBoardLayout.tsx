import ProfileSummaryCard from '../ProfileSummaryCard/profileSummaryCard';
import VisitReportLayout from '../VisitReport/visitReportLayout';
import NoteControlsLayout from '../NoteControls/noteControlsLayout';
import { HeaderInformationAtom } from './derivedAtoms';
import { useAtom } from 'jotai';
import { useMutation } from '@tanstack/react-query';
import { changeName } from '../../api/changeName';

function PatientDashBoardLayout() {
  const [headerInformation, updateHeaderInformation] = useAtom(
    HeaderInformationAtom,
  );
  const updateNameMutatuon = useMutation({
    mutationFn: changeName,
    onSuccess: () => {
      console.log('Name updated');
    },
    onError: error => {
      console.error(error);
    },
  });

  const handleUpdateName = (name: string) => {
    updateNameMutatuon.mutate({
      ReportID: headerInformation.id,
      NewName: name,
    });
  };

  return (
    <div className="flex flex-col h-full max-h-screen overflow-hidden">
      <div className="flex-shrink-0 border-b border-gray-300">
        <ProfileSummaryCard
          name={headerInformation.name}
          description={headerInformation.oneLiner}
          onChange={updateHeaderInformation} // TODO: add API call to update patient name after update
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

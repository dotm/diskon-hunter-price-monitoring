import { LocalStorageKey } from "@/utils/constants";
import { LoggedInUserData } from "@/utils/models";
import useLocalStorage from "use-local-storage";
import EditUserDataForm from "./EditUserDataForm";
import ResetPasswordForm from "./ResetPasswordForm";
import SignInForm from "./SignInForm";
import SignOutForm from "./SignOutForm";

export default function Settings() {
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)

  return (
    <div>
      {
        loggedInUserData !== undefined
        ?
        <div className="space-y-10">
          <SignOutForm/>
          <EditUserDataForm/>
        </div>
        :
        <div className="space-y-10">
          <SignInForm/>
          <ResetPasswordForm/>
        </div>
      }
    </div>
  )
}

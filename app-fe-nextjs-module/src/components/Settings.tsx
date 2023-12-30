import { LocalStorageKey } from "@/utils/constants";
import { LoggedInUserData } from "@/utils/models";
import useLocalStorage from "use-local-storage";
import ResetPasswordForm from "./ResetPasswordForm";
import SignInForm from "./SignInForm";
import SignOutForm from "./SignOutForm";

export default function Settings() {
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)

  return (
    <div className="space-y-10">
      {
        loggedInUserData === undefined
        ?
        <>
          <SignInForm/>
          <ResetPasswordForm/>
        </>
        :
        <>
          <SignOutForm/>
          {/* <EditUserDataForm/> */}
        </>
      }
    </div>
  )
}
